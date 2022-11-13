FROM public.ecr.aws/lambda/provided:al2 as builder

# install compiler
RUN yum install -y golang
RUN go env -w GOPROXY=direct

# Builder dir
RUN mkdir build
WORKDIR build

# Add dependencies
ADD go.mod go.sum ./
RUN go mod download

# Build
ADD . .
RUN go build -o /main

# Runner image
FROM public.ecr.aws/lambda/provided:al2

ENV NOTION_API_KEY ""
ENV NOTION_DATABASE_ID ""
ENV GITHUB_USER = ""

COPY --from=builder /main /main
ENTRYPOINT [ "/main" ]
