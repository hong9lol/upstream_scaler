# Python 이미지를 기반으로 설정
FROM python:3.8-slim

# 작업 디렉토리 설정
RUN mkdir /simple_app
ADD . /simple_app
WORKDIR /simple_app

# 현재 디렉토리의 cpu_load.py를 /app 디렉토리로 복사
COPY cpu_load.py /simple_app

# 컨테이너 시작 시 Python 스크립트 실행
CMD ["python", "./cpu_load.py"]
