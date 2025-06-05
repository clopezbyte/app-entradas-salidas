FROM python:3.11-slim

WORKDIR /app

COPY . /app
COPY requirements.txt /app/

RUN pip install -r requirements.txt

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]