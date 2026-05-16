import smtplib
import ssl
from datetime import datetime
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email import encoders
from typing import List, Optional

from .config import settings
from .logger import logger
from .models import Group, get_session
from .report import ReportGenerator


class EmailNotifier:
    def __init__(self):
        self.report_gen = ReportGenerator()
        self.smtp_config = settings.smtp

    def _build_message(
        self,
        subject: str,
        html_content: str,
        to_addrs: List[str],
        attachments: Optional[List[str]] = None,
    ) -> MIMEMultipart:
        msg = MIMEMultipart("alternative")
        msg["Subject"] = subject
        msg["From"] = self.smtp_config.from_addr
        msg["To"] = ", ".join(to_addrs)

        part_html = MIMEText(html_content, "html", "utf-8")
        msg.attach(part_html)

        if attachments:
            for file_path in attachments:
                try:
                    with open(file_path, "rb") as f:
                        part = MIMEBase("application", "octet-stream")
                        part.set_payload(f.read())
                    encoders.encode_base64(part)
                    import os
                    filename = os.path.basename(file_path)
                    part.add_header(
                        "Content-Disposition",
                        f"attachment; filename= {filename}",
                    )
                    msg.attach(part)
                except Exception as e:
                    logger.warning(f"Failed to attach {file_path}: {e}")

        return msg

    def _send(self, msg: MIMEMultipart, to_addrs: List[str]) -> bool:
        if not self.smtp_config.host or not self.smtp_config.username:
            logger.warning("SMTP not configured, skipping email send")
            return False

        try:
            context = ssl.create_default_context()
            if self.smtp_config.use_tls:
                server = smtplib.SMTP(self.smtp_config.host, self.smtp_config.port)
                server.starttls(context=context)
            else:
                server = smtplib.SMTP_SSL(
                    self.smtp_config.host, self.smtp_config.port, context=context
                )

            server.login(self.smtp_config.username, self.smtp_config.password)
            server.sendmail(
                self.smtp_config.from_addr,
                to_addrs,
                msg.as_string(),
            )
            server.quit()
            logger.info(f"Email sent successfully to {len(to_addrs)} recipients")
            return True
        except Exception as e:
            logger.exception(f"Failed to send email: {e}")
            return False

    def send_daily_report(
        self,
        period: str = "daily",
        to_addrs: Optional[List[str]] = None,
        group_id: Optional[int] = None,
    ) -> bool:
        to_addrs = to_addrs or self.smtp_config.to_addrs
        if not to_addrs:
            logger.warning("No recipients configured for email")
            return False

        html_content = self.report_gen.generate(period=period, fmt="html", group_id=group_id)
        md_path = self.report_gen.save_report(period=period, fmt="markdown")

        group_name = "全部"
        if group_id is not None:
            session = get_session()
            try:
                group = session.get(Group, group_id)
                group_name = group.name if group else "全部"
            finally:
                session.close()

        subject = f"📰 RSS 每日摘要 - {group_name} - {datetime.now().strftime('%Y-%m-%d')}"
        msg = self._build_message(subject, html_content, to_addrs, [md_path])

        return self._send(msg, to_addrs)

    def send_report_to_all_groups(self, period: str = "daily") -> int:
        session = get_session()
        try:
            groups = session.query(Group).all()
            sent_count = 0

            if self.send_daily_report(period=period):
                sent_count += 1

            for group in groups:
                if self.send_daily_report(period=period, group_id=group.id):
                    sent_count += 1

            return sent_count
        finally:
            session.close()
