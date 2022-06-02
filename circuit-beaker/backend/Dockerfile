FROM ruby:alpine3.9

WORKDIR /app

COPY Gemfile .
RUN bundle install && bundle clean
COPY . /app

EXPOSE 80
CMD ["ruby", "app.rb"]
