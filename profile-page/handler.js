'use strict'

const fsPromises = require('fs').promises;

module.exports = async (event, context) => {

  let index = await fsPromises.readFile('./function/static/index.html', 'utf8');

  return context
    .status(200)
    .headers({'content-type': "text/html"})
    .succeed(index)
}
