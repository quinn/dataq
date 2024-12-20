/**
 * @typedef {Object} ConfigObject
 * @property {string} method
 * @property {string} path
 * @property {string} routeFilename
 * @property {string} viewFilename
 * @property {string} funcName */

/**
 * @param {Object} options
 * @param {string} options.method
 * @param {string} options.path
 * @returns {ConfigObject} */
function config({ method, path }) {
  method = method.toUpperCase();

  let rpath = path;

  // remove first char of path if it is '/'
  if (path.startsWith("/")) {
    rpath = path.slice(1);
  }

  const parts = rpath.split("/");
  let routeFilename = parts
    .map((part) => {
      if (part.startsWith(":")) {
        return "[" + part.replace(":", "") + "]";
      }

      return part;
    })
    .join(".");

  let funcName = parts
    .filter((part) => !part.startsWith(":"))
    .map((part) => {
      return convertCase("pascal", part);
    })
    .join("");

  let reverseName = parts.filter((part) => !part.startsWith(":")).join(".");

  switch (method) {
    case "POST":
      funcName = funcName + "Create";
      routeFilename = routeFilename + "POST";
      break;
    case "DELETE":
      funcName = funcName + "Delete";
      routeFilename = routeFilename + "DELETE";
      break;
  }

  return { method, path, routeFilename, funcName, reverseName };
}

/**
 * @param {string} fileData
 * @param {ConfigObject} config
 * @returns {string} */
function addRoute(fileData, config) {
  let out = "";
  let alreadyAdded = false;
  for (const line of fileData.split("\n")) {
    if (line.includes(`e.${config.method}("${config.path}"`)) {
      alreadyAdded = true;
    }

    if (!alreadyAdded && line.endsWith("/* insert new routes here */")) {
      out += `e.${config.method}("${config.path}", routes.${config.funcName}).Name = "${config.reverseName}";\n`;
    }
    out += line + "\n";
  }

  return out;
}
