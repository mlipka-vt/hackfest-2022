# frozen_string_literal: true

require 'sinatra'
require './lib/user_services_pb'
require './lib/user_pb'
require './lib/trade_services_pb'
require './lib/trade_pb'
require 'kafka'
require 'securerandom'

kafka = Kafka.new('localhost:9092', client_id: 'ruby-client', resolve_seed_brokers: true)

userStub = UserService::Stub.new('localhost:50051', :this_channel_is_insecure)
tradeStub = TradeService::Stub.new('localhost:50052', :this_channel_is_insecure)

before do
  content_type 'application/json'
end

def to_json(request)
  request.body.rewind
  JSON.parse(request.body.read)
end

post '/login' do
  json_params = to_json(request)
  resp = userStub.login(LoginRequest.new(username: json_params['username'], password: json_params['password']))
  return resp.status
end

post '/logout' do
  json_params = to_json(request)
  resp = userStub.logout(LogoutRequest.new(username: json_params['username']))
  return resp.status
end

post '/trade/status' do
  json_params = to_json(request)
  resp = tradeStub.get_trade_status(GetTradeStatusRequest.new(id: json_params['id']))
  puts resp
  return resp.to_s
end

post '/trade' do
  json_params = to_json(request)
  id = SecureRandom.uuid

  trade = Trade.new({
                      id: id,
                      symbol: json_params['symbol'],
                      amount: json_params['amount'],
                      user: User.new(username: json_params['username'])
                    })
  bytes = Trade.encode(trade)
  kafka.deliver_message(bytes, topic: 'trades')
  return { id: id }.to_hash
end
