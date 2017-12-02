FROM scratch

ENV PORT "6060"
ENV MNG_PORT "7070"
ENV MONGO_URI "mongodb://localhost"
ENV REGISTRATION_ENABLED "true"

ADD ./dist/*

CMD [ "caloriosa-core" ]