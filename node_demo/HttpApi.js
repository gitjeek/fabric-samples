function send(promise,res) {
  promise.then((data) => {
    res.send(data);
  }).catch((err) => {
    res.send(err);
  });
}

exports.send = send;