FROM python:2

COPY requirements.txt /requirements.txt

RUN pip install -r requirements.txt

COPY server.py /server.py

CMD python /server.py /output.log stats.mailu.io.
