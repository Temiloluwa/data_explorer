FROM python:alpine3.10

ENV key=value

WORKDIR /app

COPY src/requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

COPY src .

EXPOSE 8000

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]