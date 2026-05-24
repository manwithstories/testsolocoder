import { useState, useEffect } from 'react';
import { Card, Row, Col, Tag, Empty, message, Button } from 'antd';
import { DragDropContext, Droppable, Draggable, DropResult } from 'react-beautiful-dnd';
import { CalendarOutlined } from '@ant-design/icons';
import { scheduleApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { Schedule } from '../../types';
import dayjs from 'dayjs';

export const ScheduleBoardPage = () => {
  const { user } = useAuthStore();
  const [schedules, setSchedules] = useState<Schedule[]>([]);
  const [loading, setLoading] = useState(false);
  const [weekStart, setWeekStart] = useState(dayjs().startOf('week'));

  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const params: any = {
        start_date: weekStart.format('YYYY-MM-DD'),
        end_date: weekStart.endOf('week').format('YYYY-MM-DD'),
      };
      const res = user?.role === 'temporary'
        ? await scheduleApi.getMySchedules(params)
        : await scheduleApi.getSchedules(params);
      if (res.data.code === 200) {
        setSchedules(res.data.data.data || []);
      }
    } catch (error) {
      message.error('获取排班数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSchedules();
  }, [weekStart]);

  const handleDragEnd = async (result: DropResult) => {
    if (!result.destination) return;

    const { source, destination, draggableId } = result;
    
    if (source.droppableId === destination.droppableId) return;

    const schedule = schedules.find(s => s.id === draggableId);
    if (!schedule) return;

    const newDate = dayjs(destination.droppableId);

    try {
      await scheduleApi.updateSchedule(draggableId, {
        shift_date: newDate.format('YYYY-MM-DD'),
      });
      message.success('排班已更新');
      fetchSchedules();
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败');
    }
  };

  const getWeekDays = () => {
    const days = [];
    for (let i = 0; i < 7; i++) {
      days.push(weekStart.add(i, 'day'));
    }
    return days;
  };

  const getSchedulesForDate = (date: dayjs.Dayjs) => {
    return schedules.filter(s => dayjs(s.shift_date).isSame(date, 'day'));
  };

  const weekDays = getWeekDays();

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">
          <CalendarOutlined className="mr-2" />
          排班看板 - {weekStart.format('YYYY年MM月')}
        </h2>
        <div className="flex gap-2">
          <Button onClick={() => setWeekStart(weekStart.subtract(1, 'week'))}>上一周</Button>
          <Button onClick={() => setWeekStart(dayjs().startOf('week'))}>本周</Button>
          <Button onClick={() => setWeekStart(weekStart.add(1, 'week'))}>下一周</Button>
        </div>
      </div>

      <DragDropContext onDragEnd={handleDragEnd}>
        <Row gutter={[8, 8]}>
          {weekDays.map((day) => (
            <Col key={day.format('YYYY-MM-DD')} span={24} md={12} lg={8} xl={6}>
              <Card
                title={`${day.format('MM-DD')} ${day.format('dddd')}`}
                size="small"
                className="min-h-48"
                styles={{ header: { background: day.isSame(dayjs(), 'day') ? '#e6f4ff' : undefined } }}
              >
                <Droppable droppableId={day.format('YYYY-MM-DD')}>
                  {(provided, snapshot) => (
                    <div
                      ref={provided.innerRef}
                      {...provided.droppableProps}
                      className={`min-h-32 p-2 rounded ${snapshot.isDraggingOver ? 'bg-blue-50' : ''}`}
                    >
                      {getSchedulesForDate(day).length === 0 ? (
                        <Empty description="无排班" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                      ) : (
                        getSchedulesForDate(day).map((schedule, index) => (
                          <Draggable
                            key={schedule.id}
                            draggableId={schedule.id}
                            index={index}
                            isDragDisabled={user?.role === 'temporary'}
                          >
                            {(provided, snapshot) => (
                              <div
                                ref={provided.innerRef}
                                {...provided.draggableProps}
                                {...provided.dragHandleProps}
                                className={`p-2 mb-2 rounded border ${snapshot.isDragging ? 'shadow-lg bg-white' : 'bg-gray-50'} cursor-move`}
                              >
                                <div className="flex justify-between items-center">
                                  <span className="font-medium text-sm">{schedule.temporary?.real_name}</span>
                                  <Tag
                                    color={schedule.status === 'completed' ? 'green' : schedule.status === 'in_progress' ? 'orange' : 'blue'}
                                    style={{ margin: 0 }}
                                  >
                                    {schedule.start_time}
                                  </Tag>
                                </div>
                                <div className="text-xs text-gray-500 mt-1">
                                  {schedule.job_posting?.position}
                                </div>
                              </div>
                            )}
                          </Draggable>
                        ))
                      )}
                      {provided.placeholder}
                    </div>
                  )}
                </Droppable>
              </Card>
            </Col>
          ))}
        </Row>
      </DragDropContext>
    </div>
  );
};
