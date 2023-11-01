import argparse
import os
import csv
import random
import string
import uuid
import json
import time
from datetime import datetime, timezone
from decimal import Decimal
from os.path import join, dirname
from dotenv import load_dotenv
from tqdm import tqdm
from sqlalchemy import create_engine
from sqlalchemy import select, func, update
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, Session


class Base(DeclarativeBase):
    pass


class Device(Base):
    __tablename__ = 'm_device'

    device_id: Mapped[str] = mapped_column(primary_key=True)
    tenant_id: Mapped[str]
    device_cd: Mapped[str]
    format_id: Mapped[str]


class Format(Base):
    __tablename__ = 'm_format_detail'

    format_detail_id: Mapped[str] = mapped_column(primary_key=True)
    tenant_id: Mapped[str]
    format_id: Mapped[str]


class DataStatus(Base):
    __tablename__ = 't_data_status'

    status_id: Mapped[str] = mapped_column(primary_key=True)
    tenant_id: Mapped[str]
    device_cd: Mapped[str]
    change_datetime: Mapped[str]
    status: Mapped[str]


class DataHistory(Base):
    __tablename__ = 't_data_status_history'

    status_history_id: Mapped[str] = mapped_column(primary_key=True)
    tenant_id: Mapped[str]
    device_cd: Mapped[str]
    seq: Mapped[int]
    change_datetime: Mapped[str]
    status: Mapped[str]


def create(args):
    load_dotenv(verbose=True)
    dotenv_path = join(dirname(__file__), '.env')
    load_dotenv(dotenv_path)

    ENDPOINT = os.environ.get("ENDPOINT")
    PORT = os.environ.get("PORT")
    DBNAME = os.environ.get("DBNAME")
    USERNAME = os.environ.get("USERNAME")
    PASSWORD = os.environ.get("PASSWORD")

    CONNECT_STR = '{}://{}:{}@{}:{}/{}'.format(
        'postgresql', USERNAME, PASSWORD, ENDPOINT, PORT, DBNAME)

    engine = create_engine(CONNECT_STR, echo=False)
    with Session(engine) as session:
        device = session.scalar(select(Device)
                                .where(Device.device_cd == args.id))

        formats = session.scalar(select(func.count())
                                 .select_from(Format)
                                 .where(Format.format_id == device.format_id))

        status = session.scalar(select(func.count())
                                .select_from(DataStatus)
                                .where(DataStatus.device_cd == args.id)
                                .where(DataStatus.tenant_id == device.tenant_id))

        dummy_data = []

        with open('./' + args.file, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            for row in reader:
                if formats == len(row):
                    random_list = [random.choice(
                        string.hexdigits) for n in range(14)]
                    data_id = "".join(random_list)
                    data = {'status0': data_id.upper()}

                    for i, col in enumerate(row):
                        key = i + 1
                        data['status' + str(key)] = col

                dummy_data.append(data)

        for i, data in enumerate(tqdm(dummy_data)):
            data = json.dumps(data)
            now = str(datetime.now(timezone.utc))
            status_uuid = str(uuid.uuid4())
            seq = session.scalar(select(func.max(DataHistory.seq))
                                 .select_from(DataHistory)
                                 .where(DataHistory.device_cd == args.id))

            if seq == None:
                seq = 0

            seq += 1

            session.add(DataHistory(
                status_history_id=status_uuid,
                tenant_id=device.tenant_id,
                device_cd=args.id,
                seq=seq,
                change_datetime=now,
                status=data
            ))

            if i + 1 == len(dummy_data):
                if status == 0:
                    session.add(DataStatus(
                        status_id=status_uuid,
                        tenant_id=device.tenant_id,
                        device_cd=args.id,
                        change_datetime=now,
                        status=data
                    ))
                else:
                    session.execute(update(DataStatus)
                                    .where(DataStatus.device_cd == args.id)
                                    .where(DataStatus.tenant_id == device.tenant_id)
                                    .values(change_datetime=now, status=data)
                                    )

            time.sleep(0.01)

        session.commit()


def main():
    parser = argparse.ArgumentParser(description='')
    subparser = parser.add_subparsers()

    parser_device = subparser.add_parser(
        'dummy', help='Script for dummy data input')
    parser_device.add_argument(
        '--id', '-i', required=True, help='Required: Please input device code')
    parser_device.add_argument(
        '--file', '-f', default='dummy.csv', help='Optional: Please select input data file')
    parser_device.set_defaults(handler=create)

    args = parser.parse_args()
    if hasattr(args, 'handler'):
        args.handler(args)
    else:
        parser.print_help()


if __name__ == '__main__':
    main()
