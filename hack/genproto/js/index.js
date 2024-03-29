const fs = require('fs');
const protoLoader = require('@grpc/proto-loader');

const PROTO_GEN = `${__dirname}/../../../api/go/gen`;

const GEN_GW = new RegExp(/\.(pb.gw.go)/g);
const GEN_PB = new RegExp(/\.(pb.go)/g);

function appendTrailingSlash(str = '') {
  let _str = str;
  if (str.charAt(str.length - 1) !== '/') {
    _str = `${str}/`;
  }
  return _str;
}

// Extracted from grpc-node
// https://github.com/grpc/grpc-node
function loadPackageDefinition(packageDef) {
  const result = {};
  for (const serviceFqn in packageDef) {
    const service = packageDef[serviceFqn];
    const nameComponents = serviceFqn.split('.');
    const serviceName = nameComponents[nameComponents.length - 1];
    let current = result;
    // We are only interested in implemented services
    // Apparently grpc package `loadPackageDefinition` uses the same logic
    // to check for services. I could have misinterpret it.
    if (!service.hasOwnProperty('format')) {
      current.serviceName = serviceName;
      current.methods = Object.keys(service).map(k => ({
        methodName: k,
        path: service[k].path
      }));
      current.packageSvc = serviceFqn;
      current.package = nameComponents
        .slice(0, nameComponents.length - 1)
        .join('.');
      current.currentPackage = nameComponents[1];
      current.originalPackage = nameComponents
        .slice(0, nameComponents.length - 1)
        .join('_');
    }
  }
  return Object.keys(result).length > 0 ? result : null;
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
  const packageDefinition = protoLoader.loadSync(filePath, options);

  return loadPackageDefinition(packageDefinition);
}

/**
 * Loads proto files
 * @param {Array} refProtos adds proto by reference
 * @param {String} filePath Dir
 * @param {Array} relativeInclude Directory has to be relative to where it is being loaded from
 */
function loadProtoDir({ refProtos = [], path = [], filePath, includeDirs }) {
  if (!filePath) filePath = path;
  fs.readdirSync(filePath).forEach(fileName => {
    const fullPath = appendTrailingSlash(filePath) + fileName;
    if (!fs.statSync(fullPath).isFile()) {
      // Folder
      loadProtoDir({ refProtos, path, filePath: fullPath, includeDirs });
    } else if (fileName.match(/\.proto$/) && !filePath.match(/third_party/)) {
      // exclude third party
      const proto =
        !includeDirs || !includeDirs.length
          ? loadProto(fullPath)
          : loadProto(fullPath, includeDirs);
      if (!proto) {
        return;
      }
      const parent = includeDirs.find((i) => path.includes(i));
      const f = path.slice(parent.length);
      const fp = fullPath.slice(path.length);

      proto.path = `${f}`
      proto.filePath = fp;
      refProtos.push(proto);
    }
  });
}

function loadMap(mapFile) {
  const protos = [];
  if (!Array.isArray(mapFile)) {
    throw new Error("Map File has to an array");
  }
  mapFile.forEach((m) => {
    m.includes = m.includes.map((i) => `${__dirname}/../../../${i}`)
    m.protos.forEach((p) =>
      loadProtoDir({
        refProtos: protos,
        path: `${__dirname}/../../../${p}`,
        includeDirs: m.includes
      }))
  })
  return protos;
}

// Checks if protoc gateway has been generated
// for the file.
function ensureGatewayGenerated(initialPath, filePath, object = {}) {
  if (!filePath) filePath = initialPath;
  fs.readdirSync(filePath).forEach(fileName => {
    const fullPath = appendTrailingSlash(filePath) + fileName;
    // Folder
    if (!fs.statSync(fullPath).isFile()) {
      ensureGatewayGenerated(initialPath, fullPath, object);
    } else if (GEN_GW.test(fileName)) {
      const p = initialPath.split('/');
      let f = filePath.split('/');
      f = f.slice(p.length).join('/');
      object[`/${f}`] = true;
    }
  });
}

function writeToFile(obj) {
  const str = JSON.stringify(obj);
  fs.writeFile(`${__dirname}/../descriptor.json`, str, 'utf8', err => {
    if (err) {
      console.log('An error occured while writing JSON Object to File.\n', err);
      return;
    }
  });
}

function getArgs() {
  const args = process.argv.slice(2);
  switch (args[0]) {
    case "--map": case "-m":
      const data = args[1];
      return JSON.parse(data);
  }
}

function main() {
  console.log('index.js: Generating Descriptors');
  const mapFile = getArgs();
  const condition = {};
  ensureGatewayGenerated(PROTO_GEN, null, condition);
  const protos = loadMap(mapFile);
  // loadPaths({ protos, paths: PROTOS, includeDirs: INCLUDES });
  // Match based on file path
  writeToFile(protos.filter(p => condition[p.path] !== undefined));
  console.log('index.js: JSON file has been saved.');
}
main();
