import {
  fetchGraphData,
  fetchGraphDefinition,
  normalizeData,
} from "./chart-data.js";

/**
 * Adds a script to the document dynamically.
 * @param {string} src - The source URL of the script to be added.
 * @returns {Promise<void>} - A promise that resolves when the script is loaded successfully, or rejects if there is an error.
 */
const addScript = async (src) =>
  new Promise((resolve, reject) => {
    const existingEl = document.querySelector(`script[src="${src}"]`);
    if (existingEl) {
      existingEl.addEventListener("load", resolve);
      return; // Script already exists, no need to add it again
    }

    const el = document.createElement("script");
    el.src = src;
    el.crossorigin = "";
    el.addEventListener("load", resolve);
    el.addEventListener("error", reject);
    document.body.append(el);
  });

/**
 * A template literal tag 'html' that creates an HTML <template> element from the contents of the string,
 *
 * - It's efficient: The browser doesn't render the content of a <template> element until it's used.
 * - It's reusable: You can clone the template's content multiple times.
 *  - It's safe: The content is parsed as HTML, so you don't need to worry about XSS attacks from the  interpolated values.
 * @param {Array<string>} strings - An array of strings representing the HTML template.
 * @param {...any} values - The dynamic values to be inserted into the HTML template.
 * @returns {HTMLTemplateElement} The template element created from the HTML string.
 */
function html(strings, ...values) {
  // Combine the strings and values
  const rawHTML = strings.reduce((result, string, i) => {
    return result + string + (values[i] || "");
  }, "");

  // Create a template element
  const template = document.createElement("template");

  // Set the innerHTML of the template
  template.innerHTML = rawHTML.trim();

  // Return the template element
  return template;
}
const URLs = {
  echarts: "https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js",
};

class WikiChart extends HTMLElement {
  constructor() {
    super();

    this.connected = false;
    this.rendered = false;
    this.wikichart = null;
    // check for a Declarative Shadow Root.
    const shadow = this.internals?.shadowRoot;
    if (!shadow) {
      this.attachShadow({ mode: "open" });
      if (this.constructor.template) {
        this.shadowRoot.setHTMLUnsafe(this.constructor.template.innerHTML);
      }
    }
    if (this.constructor.stylesheetURL) {
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = this.constructor.stylesheetURL;
      this.shadowRoot.appendChild(link);
    }
    this.initializeProperties();
  }

  initializeProperties() {
    const props = this.constructor.properties;
    for (const [name, config] of Object.entries(props)) {
      if (!Object.prototype.hasOwnProperty.call(this, name)) {
        const value = this.getAttribute(name) || config.default;
        if (config.options && !config.options.includes(value)) {
          console.warn(
            `Invalid value ${value} for property ${name}. Valid options are ${config.options}`,
          );
          continue;
        }
        this[name] = this.convertValue(value, config.type);
      }
    }
  }

  convertValue(value, type) {
    if (value === null) {
      return value;
    }
    switch (type) {
      case String:
        return value;
      case Number:
        return Number(value);
      case Boolean:
        return value !== null && value !== "false";
      case Array:
      case Object:
        if (typeof value === "string") {
          return JSON.parse(value);
        }
        return value;

      default:
        return value;
    }
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue !== newValue) {
      const props = this.constructor.properties;
      if (name in props) {
        this[name] = this.convertValue(newValue, props[name].type);
        this.propertyChangedCallback(name, oldValue, this[name]);
      }
    }
  }

  propertyChangedCallback(name, oldValue, newValue) {
    console.log(`Property ${name} changed from ${oldValue} to ${newValue}`);
    // This method can be overridden in subclasses to react to property changes
    this.render();
  }

  connectedCallback() {
    this.connected = true;
    if (this.chart) {
      this.render();
    }
  }

  disconnectedCallback() {
    if (super.disconnectedCallback) {
      super.disconnectedCallback();
      this.connected = false;
      this.wikichart?.dispose();
    }
  }

  static get observedAttributes() {
    return Object.keys(this.properties);
  }

  static get template() {
    return html`
      <div
        class="wiki-chart"
        style="width: 100%;height:100%;min-height:500px;"
      ></div>
    `;
  }

  static get properties() {
    return {
      chart: {
        type: String,
      },
      "data-chart": {
        type: String,
      },
      theme: {
        type: String,
        options: ["light", "dark"],
        default: "light",
      },
    };
  }

  async render() {
    if (this["data-chart"]) {
      await this.renderGraph(
        JSON.parse(decodeURIComponent(this["data-chart"])),
      );
    } else if (this.chart) {
      const url = new URL(this.chart);
      this.article = url.pathname.split("/wiki/").pop();
      this.hostname = url.hostname;
      this.title = this.article.split(".")[0].replace(/_/g, " ");
      const graphDefinition = await fetchGraphDefinition(
        this.hostname,
        this.article,
      );
      const graphData = await fetchGraphData(
        this.hostname,
        `Data:${graphDefinition.source}`,
      );
      const spec = {
        locale: "en",
        data: normalizeData(graphData),
        definition: normalizeData(graphDefinition),
      };
      await this.renderGraph(GetEchartOptions(spec));
    }
  }

  async renderGraph(spec) {
    if (!window.echarts) {
      await addScript(URLs.echarts);
    }
    this.wikichart?.dispose();
    // Display the chart using the configuration items and data just specified.
    this.wikichart = window.echarts.init(
      this.shadowRoot.querySelector(".wiki-chart"),
      this.theme,
      { renderer: "svg" },
    );

    this.wikichart.setOption(spec);
    window.addEventListener("resize", this.wikichart.resize);
  }
}

customElements.define("wiki-chart", WikiChart);
