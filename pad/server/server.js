var messages = require('./changesetproto/changeset_pb');
var services = require('./changesetproto/changeset_grpc_pb');

var grpc = require('@grpc/grpc-js');

function Generate(call, callback) {
  var request = call.request;
  console.log(request);
  oldtext = request.getOldtext();
  newtext = request.getNewtext();
  attribs = request.getAttribs();

  chset = GenerateChangeset(oldtext, newtext, attribs);

  var reply = new messages.GenerateReply();
  reply.setChangeset(chset);

  callback(null, reply);
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
  var server = new grpc.Server();
  server.addService(services.ChangesetService, {generate: Generate});
  server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {
    server.start();
  });
}

main();


var Changeset = require("./etherpad-lite/src/static/js/Changeset");

function GenerateChangeset(oldtext, newtext, attribs) {
  // init the changeset builder
  var builder = Changeset.builder(oldtext.length);

  // find the longest common prefix
  var commonLength = 0;
  while (commonLength < oldtext.length && commonLength < newtext.length && oldtext[commonLength] === newtext[commonLength]) {
      commonLength++;
  }

  // keep common prefix
  builder.keep(commonLength);

  // // remove the remaining of the old text
  // if (commonLength < oldtext.length) {
  //     builder.remove(oldtext.length - commonLength);
  // }

  // add the remaining of the new text
  if (commonLength < newtext.length) {
      builder.insert(newtext.substring(commonLength));
  }

  // generate the changeset
  return builder.toString();
}

// function GenerateChangeset(oldtext, newtext, attribs) {
//   // init the changeset builder
//   var builder = Changeset.builder(oldtext.length);

//   // insert the new text
//   builder.insert(newtext);

//   // generate the changeset
//   return builder.toString();
// }

