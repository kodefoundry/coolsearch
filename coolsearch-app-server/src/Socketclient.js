import openSocket from 'socket.io-client';
const  socket = openSocket('http://localhost:4050');

function subscribeToSearchResult(cb) {
  socket.on('broadcast', uuid => cb(null, uuid));
}
export { subscribeToSearchResult };