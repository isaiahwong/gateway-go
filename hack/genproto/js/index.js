const fs = require('fs');
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

const PROTOS = `${__dirname}/../../../proto/api`
const INCLUDES = [
  `${__dirname}/../../../proto/api`,
  `${__dirname}/../../../proto/third_party/googleapis`,
];

function MethodDesc() {
  this.methodName = "";
  this.path = "";
}

function ServiceDesc() {
  this.serviceName = "";
  this.originalName = "";
  this.package = "";
  this.methods = [];
}

function appendTrailingSlash(str = '') {
  let _str = str;
  if (str.charAt(str.length - 1) !== '/') {
    _str = `${str}/`;
  }
  return _str;
}

function loadProto(filePath, include) {
  const options = {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  };
  if (!fs.statSync(filePath).isFile()) {
    throw new InternalServerError('proto file not found');
  }
  if (Array.isArray(include) && include.length) {
    options.includeDirs = [...include];
  }
  const packageDefinition = protoLoader.loadSync(
    filePath,
    options
  );
  return grpc.loadPackageDefinition(packageDefinition);
}

/**
 * Loads proto files
 * @param {Array} protos adds proto by reference
 * @param {String} filePath Dir
 * @param {Array} relativeInclude Directory has to be relative to where it is being loaded from
 */
function loadProtos(protos = [], filePath, relativeInclude) {
  fs
    .readdirSync(filePath)
    .forEach((fileName) => {
      const fullPath = appendTrailingSlash(filePath) + fileName;
      if (!fs.statSync(fullPath).isFile()) { // Folder
        loadProtos(protos, fullPath, relativeInclude);
      }
      else if (fileName.match(/\.proto$/) && !filePath.match(/third_party/)) { // exclude third party
        const proto = (!relativeInclude || !relativeInclude.length)
          ? loadProto(fullPath)
          : loadProto(fullPath, relativeInclude);
        protos.push(proto);
      }
    });
}

function load() {
  const protos = []
  loadProtos(protos, PROTOS, INCLUDES)
  return protos.filter((p) => {
    const protoPkg = Object.keys(p)[0];
    if (!protoPkg) {
      return false;
    }
    const service = Object.keys(p[protoPkg])[0];
    if (!service) {
      return false;
    }
    return !!p[protoPkg][service].service;
  });
}

function serialize(protos, _services = []) {
  if (!protos || !protos.length) {
    logger.warn('No protos supplied');
    return;
  }

  protos.forEach((proto) => {
    const _package = Object.keys(proto)[0];
    if (!_package) {
      return;
    }
    const services = proto[_package];
    if (!services) {
      return;
    }
    Object.keys(services).forEach((key) => {
      if (typeof services[key] !== 'function') return;
      const d = new ServiceDesc();
      d.serviceName = `${_package}.${key}`;
      d.originalName = key
      d.package = _package
      d.methods = Object.keys(services[key].service)
        .map((k) => {
          const m = new MethodDesc();
          m.methodName = k;
          m.path = services[key].service[k].path;
          return m;
        });
      _services.push(d);
    });
  });
}

function writeToFile(obj) {
  const str = JSON.stringify(obj);
  fs.writeFile(`${__dirname}/../descriptor.json`, str, 'utf8', (err) => {
    if (err) {
        console.log('An error occured while writing JSON Object to File.\n', err);
        return
    }
    console.log('JSON file has been saved.');
});
}

function main() {
  const _services = [];
  serialize(load(), _services);
  writeToFile(_services);
}
main();
