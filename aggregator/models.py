from sqlalchemy import Column, String, Integer
from sqlalchemy.types import TIMESTAMP
from database import Base

class ClientSession(Base):
    __tablename__ = "client_sessions"

    environment = Column(String(255), nullable=False)
    realm_id = Column(String(255), nullable=False)
    client_id = Column(String(255), nullable=False)
    count = Column(Integer, nullable=False)
    date = Column(TIMESTAMP(timezone=True), nullable=False)
    id = Column(Integer, primary_key=True)

class ClientEvent(Base):
    __tablename__ = "client_events"

    environment = Column(String(255), nullable=False, primary_key=True)
    realm_id = Column(String(255), nullable=False, primary_key=True)
    client_id = Column(String(255), nullable=False, primary_key=True)
    event_type = Column(String(255), nullable=False, primary_key=True)
    date = Column(TIMESTAMP(timezone=True), nullable=False, primary_key=True)
    count = Column(Integer, nullable=False)
