var kafka = require('kafka-node');
var Client = kafka.KafkaClient;
const io = require('socket.io')();

var kclient = new Client({ kafkaHost: 'localhost:9092' });

io.on('connection', (socket) => {
  socket.on('subscribeToSearchResult', (inteval) => {
    console.log('A client is subscribing to socket connection at interval: ', interval);
  });

  const topics = [{topic: "search-available"}];

  const options = {
      autoCommit: true,
      autoCommitIntervalMs: 5000,
      fetchMaxWaitMs: 1000,
      fetchMaxBytes: 1024 * 1024,
      fromOffset: true
  };

  const consumer = new kafka.Consumer(kclient, topics, options);

  consumer.on("message", function(message) {
    console.log("New Message received: ", message.value);
    io.sockets.emit('broadcast',message.value);
  });

  consumer.on("error", function(err) {
      console.log("Error in consumer: ", err);
  });

  process.on("SIGINT", function() {
      consumer.close(true, function() {
          console.log("Closing kafka connection...")
      });
      io.close(true, function(){
          console.log("Closing socket connection...")
      })
  });
});

const port = 4050;
io.listen(port);
console.log('Listening for socket connectionson on port:', port);