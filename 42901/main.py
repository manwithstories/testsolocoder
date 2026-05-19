from fastapi import FastAPI
from app.database import init_db
from app.routers import family, records, statistics, budget

app = FastAPI(
    title="家庭记账 API",
    description="家庭记账后端服务，支持成员管理、记账记录、统计分析和预算管理",
    version="1.0.0"
)


@app.on_event("startup")
def startup_event():
    init_db()


@app.get("/", summary="健康检查")
def root():
    return {"message": "家庭记账 API 服务运行正常"}


app.include_router(family.router)
app.include_router(records.router)
app.include_router(statistics.router)
app.include_router(budget.router)
