import { createServer } from "http";
import "./assets/wasm-exec.js";
import { createReadStream, readFileSync, existsSync } from "fs";
import {
  fetchGraphData,
  fetchGraphDefinition,
  normalizeData,
} from "./assets/chart-data.js";
import URL from "url";
import { init, registerTheme } from "./assets/echarts.esm.js";

const PORT = process.env.PORT || 8080;
const categoricalTheme = {
  backgroundColor: "var(--background-color-base, #fff )",
  color: [
    "#4b77d6",
    "#eeb533",
    "#fd7865",
    "#80cdb3",
    "#269f4b",
    "#b0c1f0",
    "#9182c2",
    "#d9b4cd",
    "#b0832b",
    "#a2a9b1",
  ],
};

function initWasm(path) {
  const go = new globalThis.Go();

  return new Promise((resolve, reject) => {
    const wasmBuffer = readFileSync(path);
    WebAssembly.instantiate(wasmBuffer, go.importObject)
      .then((result) => {
        go.run(result.instance);
        resolve(result.instance);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

async function getSpec(hostname, article) {
  const graphDefinition = await fetchGraphDefinition(
    hostname + ".wikimedia.org",
    article
  );
  const graphData = await fetchGraphData(
    hostname + ".wikimedia.org",
    `Data:${graphDefinition.source}`
  );

  const echartOption = GetEchartOptions({
    locale: "en",
    data: normalizeData(graphData),
    definition: normalizeData(graphDefinition),
  });

  return echartOption;
}

async function render(spec) {
  registerTheme("categorical", categoricalTheme);
  const chart = init(null, "categorical", {
    renderer: "svg",
    ssr: true,
    width: 1000,
    height: 800,
  });

  chart.setOption(spec);
  const rawSvg = chart.renderToSVGString();
  chart.dispose();
  return rawSvg;
}

const server = createServer(async (req, res) => {
  const { query } = URL.parse(req.url, true);
  const route = req.url;

  if (route === "/") {
    createReadStream("index.html").pipe(res);
    return;
  }

  const chartmatch = route.match(/^\/chart\/(.+)\/(.+)/);
  const svgmatch = route.match(/^\/svg\/(.+)\/(.+)/);
  const specmatch = route.match(/^\/spec\/(.+)\/(.+)/);
  const assetsmatch = route.match(/^\/assets|build\/(.+)/);
  if (req.method === "GET") {
    if (assetsmatch) {
      try {
        if (!existsSync("." + route)) {
          throw new Error("File not found");
        }
        const mimeType = {
          ".js": "text/javascript",
          ".css": "text/css",
          ".wasm": "application/wasm",
        }[route.match(/\.\w+$/)[0]];
        res.setHeader("Content-Type", mimeType);
        createReadStream("." + route).pipe(res);
      } catch (error) {
        res.statusCode = 404;
        res.end("Asset not found!" + route);
      }
      return;
    }

    if (svgmatch) {
      const domain = svgmatch[1];
      const chartname = svgmatch[2];
      const spec = await getSpec(domain, chartname);
      const svg = await render(spec);
      res.statusCode = 200;
      res.setHeader("Content-Type", "image/svg+xml");
      res.end(svg);
      return;
    }

    if (specmatch) {
      const domain = specmatch[1];
      const chartname = specmatch[2];
      const spec = await getSpec(domain, chartname);
      res.statusCode = 200;
      res.setHeader("Content-Type", "application/json");
      res.end(JSON.stringify(spec));
      return;
    }

    if (chartmatch) {
      const domain = chartmatch[1];
      const chartname = chartmatch[2];
      const spec = await getSpec(domain, chartname);
      const svg = await render(spec);
      res.statusCode = 200;
      res.setHeader("Content-Type", "text/html");
      res.write(`<wiki-chart data-chart="${encodeURIComponent(JSON.stringify(spec))}" >`);
      res.write(`<template shadowrootmode="open">`);
      res.write(
        `<img src="/svg/${domain}/${chartname}" alt="${chartname}" width="100%" height="100%"/>`
      );
      res.write(`</template>`);
      res.write(`</wiki-chart>`);

      res.write(`<script src="/assets/chart-data.js" type="module"></script>`);
      res.write(`<script src="/assets/wiki-chart.js" type="module"></script>`);
      res.end(`<script src="/assets/echarts.esm.js" type="module"></script>`);
      return;
    }

    res.statusCode = 404;
    res.end("Page not found: " + route);
  }
});

initWasm("./build/chartadapter.wasm").then(() => {
  console.log("Wasm module loaded");
  server.listen(PORT, () => {
    console.log(`Starting server at port ${PORT}`);
  });
});
