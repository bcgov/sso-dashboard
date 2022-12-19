from sqlalchemy import Column, String, Integer
from sqlalchemy.types import TIMESTAMP
from database import Base

class ClientEvent(Base):
    __tablename__ = "client_events"

    environment = Column(String(255), nullable=False, primary_key=True)
    realm_id = Column(String(255), nullable=False, primary_key=True)
    client_id = Column(String(255), nullable=False, primary_key=True)
    event_type = Column(String(255), nullable=False, primary_key=True)
    date = Column(TIMESTAMP(timezone=True), nullable=False, primary_key=True)
    count = Column(Integer, nullable=False)
