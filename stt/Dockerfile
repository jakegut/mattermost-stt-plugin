FROM golang:1.17

WORKDIR /usr/src/coqui

RUN apt-get update && apt-get install -y python3 python3-pip ffmpeg libc6-dev
RUN pip3 install --extra-index-url https://google-coral.github.io/py-repo/ tflite_runtime
RUN wget https://github.com/coqui-ai/STT/releases/download/v1.2.0/native_client.tflite.Linux.tar.xz
RUN tar xf native_client.tflite.Linux.tar.xz
RUN wget https://coqui.gateway.scarf.sh/english/coqui/v1.0.0-large-vocab/model.tflite 
RUN wget https://coqui.gateway.scarf.sh/english/coqui/v1.0.0-large-vocab/large_vocabulary.scorer

ENV CGO_LDFLAGS="-L/usr/src/coqui"
ENV CGO_CXXFLAGS="-I/usr/src/coqui"
ENV LD_LIBRARY_PATH="/usr/src/coqui/:$LD_LIBRARY_PATH"
ENV TRANSCRIBER_MODEL_PATH=/usr/src/coqui/model.tflite
ENV SCORER_MODEL_PATH=/usr/src/coqui/large_vocabulary.scorer

WORKDIR /usr/app/stt

COPY . .

RUN go get .
RUN go build .

ENTRYPOINT [ "./stt" ]


