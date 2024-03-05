import os

# 디렉토리 경로 설정
directory_path = "./temp/"

# 디렉토리 내의 모든 파일과 디렉토리 목록을 가져옴
files_and_directories = os.listdir(directory_path)

# 파일명만 필터링 (디렉토리는 제외)
files = [
    f for f in files_and_directories if os.path.isfile(os.path.join(directory_path, f))
]

sorted_files = sorted(files)
idx = 0
for f in sorted_files:
    file_name = directory_path + f
    idx += 1
    # 파일을 열고 3번째 줄을 읽어 출력합니다.
    # print(idx)
    if idx % 2 != 0:
        continue
    with open(file_name, "r") as file:
        # 각 줄을 순회하면서 3번째 줄을 찾습니다.
        for i, line in enumerate(file, start=1):
            if i == 3:
                print(line.strip())  # 줄바꿈 문자를 제거하고 출력
                break  # 3번째 줄을 찾았으므로 루프를 종료합니다.
