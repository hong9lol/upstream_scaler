#!/bin/bash

# 이름의 최소 길이와 최대 길이 설정
MIN_LENGTH=3
MAX_LENGTH=8

# 알파벳 문자 배열
ALPHABET=("a" "b" "c" "d" "e" "f" "g" "h" "i" "j" "k" "l" "m" "n" "o" "p" "q" "r" "s" "t" "u" "v" "w" "x" "y" "z")

# 사용자로부터 생성할 이름의 수를 입력받음
echo "Enter the number of names you want to generate:"
read COUNT

# 입력받은 숫자만큼 랜덤 이름 생성
for (( n=1; n<=COUNT; n++ ))
do
  # 랜덤 이름 길이 생성
  NAME_LENGTH=$(($MIN_LENGTH + $RANDOM % ($MAX_LENGTH - $MIN_LENGTH + 1)))

  # 빈 문자열로 이름 초기화
  NAME=""

  # 랜덤 이름 생성
  for (( i=0; i<$NAME_LENGTH; i++ )); do
    # 랜덤 인덱스를 선택하여 문자 추가
    RANDOM_INDEX=$(($RANDOM % ${#ALPHABET[@]}))
    NAME+=${ALPHABET[$RANDOM_INDEX]}
  done

  # 첫 글자를 대문자로 변환
  NAME=$(echo "$NAME" | sed 's/^./\u&/')

  # 생성된 이름 출력
  kubectl create deployment $NAME --image=hong9lol/cpu-load-image:0.1
  echo "Random Name $n: $NAME"
done
