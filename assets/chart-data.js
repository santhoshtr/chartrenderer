export async function fetchGraphDefinition(hostname, article) {
  try {
    const apiURL = `https://${hostname}/w/api.php?action=query&prop=revisions&rvprop=content&titles=${encodeURIComponent(article)}&format=json&formatversion=2&origin=*`;
    const response = await fetch(apiURL);
    if (!response.ok) {
      throw new Error("Network response was not ok for " + apiURL);
    }
    const data = await response.json();
    return JSON.parse(data.query.pages[0]?.revisions[0]?.content);
  } catch (error) {
    console.error("Fetch error:", error);
  }
}

export async function fetchGraphData(hostname, article) {
  try {
    const apiURL = `https://${hostname}/w/api.php?action=query&prop=revisions&rvprop=content&titles=${encodeURIComponent(article)}&format=json&formatversion=2&origin=*`;
    const response = await fetch(apiURL);
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    const data = await response.json();
    return JSON.parse(data.query.pages[0]?.revisions[0]?.content);
  } catch (error) {
    console.error("Fetch error:", error);
  }
}

export function normalizeData(data) {
  // If license is string, put it under code key
  if (typeof data.license === "string") {
    data.license = { code: data.license };
  }
  if (typeof data.title === "string") {
    data.title = { en: data.title };
  }

  if (typeof data.xAxis?.title === "string") {
    data.xAxis.title = { en: data.xAxis.title };
  }
  if (typeof data.yAxis?.title === "string") {
    data.yAxis.title = { en: data.yAxis.title };
  }
  return data;
}
