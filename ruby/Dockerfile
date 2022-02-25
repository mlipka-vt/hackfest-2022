FROM ruby:latest

RUN bundle config --global frozen 1

WORKDIR /usr/src/app

COPY Gemfile Gemfile.lock ./
RUN bundle install

EXPOSE 4567

COPY . .
CMD ["bundle", "exec", "rackup", "--host", "0.0.0.0", "-p", "4567"]
