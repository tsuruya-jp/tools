import argparse
import os
import csv
from datetime import datetime
from decimal import Decimal
from os.path import join, dirname
from dotenv import load_dotenv
from tqdm import tqdm
from sqlalchemy import create_engine
from sqlalchemy import insert
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, Session


class Base(DeclarativeBase):
    pass


class Invoice(Base):
    __tablename__ = 't_invoice'

    invoice_seq: Mapped[int] = mapped_column(
        primary_key=True, autoincrement=True)
    order_id: Mapped[str]
    billing_date: Mapped[datetime]
    user_name: Mapped[str]
    amount_tax: Mapped[Decimal]
    contract_devices: Mapped[str]
    create_user_seq: Mapped[int]
    create_datetime: Mapped[datetime]
    company_name: Mapped[str]


class Receipt(Base):
    __tablename__ = 't_receipt'

    receipt_seq: Mapped[int] = mapped_column(
        primary_key=True, autoincrement=True)
    order_id: Mapped[str]
    payment_date: Mapped[datetime]
    user_name: Mapped[str]
    amount_tax: Mapped[Decimal]
    create_user_seq: Mapped[int]
    create_datetime: Mapped[datetime]
    company_name: Mapped[str]


class Quotation(Base):
    __tablename__ = 't_quotation'

    quotation_seq: Mapped[int] = mapped_column(
        primary_key=True, autoincrement=True)
    qutation_no: Mapped[str]
    qutation_date: Mapped[datetime]
    user_name: Mapped[str]
    amount_tax: Mapped[Decimal]
    contract_devices: Mapped[str]
    create_user_seq: Mapped[int]
    create_datetime: Mapped[datetime]
    company_name: Mapped[str]


def db_conf():
    load_dotenv(verbose=True)
    dotenv_path = join(dirname(__file__), '.env')
    load_dotenv(dotenv_path)

    ENDPOINT = os.environ.get("ENDPOINT")
    PORT = os.environ.get("PORT")
    DBNAME = os.environ.get("DBNAME")
    USERNAME = os.environ.get("USERNAME")
    PASSWORD = os.environ.get("PASSWORD")

    return '{}://{}:{}@{}:{}/{}'.format(
        'mysql+mysqlconnector', USERNAME, PASSWORD, ENDPOINT, PORT, DBNAME)


def invoice(args):
    con = db_conf()
    engine = create_engine(con, echo=False)
    with Session(engine) as session:
        dummy_data = []

        with open('./' + args.file, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            now = datetime.now()
            for row in tqdm(reader):
                d = {
                    "order_id": row[0],
                    "billing_date": row[1],
                    "user_name": row[2],
                    "amount_tax": row[3],
                    "contract_devices": row[4],
                    "create_user_seq": row[5],
                    "create_datetime": now,
                    "company_name": row[6]
                }
                dummy_data.append(d)

        session.execute(insert(Invoice), dummy_data)

        session.commit()


def receipt(args):
    con = db_conf()
    engine = create_engine(con, echo=False)
    with Session(engine) as session:
        dummy_data = []

        with open('./' + args.file, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            now = datetime.now()
            for row in tqdm(reader):
                d = {
                    "order_id": row[0],
                    "payment_date": row[1],
                    "user_name": row[2],
                    "amount_tax": row[3],
                    "create_user_seq": row[5],
                    "create_datetime": now,
                    "company_name": row[6]
                }
                dummy_data.append(d)

        session.execute(insert(Receipt), dummy_data)

        session.commit()


def quotation(args):
    con = db_conf()
    engine = create_engine(con, echo=False)
    with Session(engine) as session:
        dummy_data = []

        with open('./' + args.file, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            now = datetime.now()
            for row in tqdm(reader):
                d = {
                    "qutation_no": row[0],
                    "qutation_date": row[1],
                    "user_name": row[2],
                    "amount_tax": row[3],
                    "contract_devices": row[4],
                    "create_user_seq": row[5],
                    "create_datetime": now,
                    "company_name": row[6]
                }
                dummy_data.append(d)

        session.execute(insert(Quotation), dummy_data)

        session.commit()


def main():
    parser = argparse.ArgumentParser(description='')
    subparser = parser.add_subparsers()

    parser_invoice = subparser.add_parser(
        'invoice', help='Script for invoice dummy data input')
    parser_invoice.add_argument(
        '--file', '-f', default='dummy.csv', help='Optional: Please select input data file')
    parser_invoice.set_defaults(handler=invoice)

    parser_receipt = subparser.add_parser(
        'receipt', help='Script for receipt dummy data input')
    parser_receipt.add_argument(
        '--file', '-f', default='dummy.csv', help='Optional: Please select input data file')
    parser_receipt.set_defaults(handler=receipt)

    parser_quotation = subparser.add_parser(
        'quotation', help='Script for quotation dummy data input')
    parser_quotation.add_argument(
        '--file', '-f', default='dummy.csv', help='Optional: Please select input data file')
    parser_quotation.set_defaults(handler=quotation)

    args = parser.parse_args()
    if hasattr(args, 'handler'):
        args.handler(args)
    else:
        parser.print_help()


if __name__ == '__main__':
    main()
