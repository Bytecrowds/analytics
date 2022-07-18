import Router from "@tsndr/cloudflare-worker-router";
import { Redis } from "@upstash/redis/cloudflare";

const router = new Router();

// Because cloudflare workers expose the env on fetch, we need to manually set the vars.
const bytecrowds = new Redis({
  url: "https://eu2-upright-dane-30471.upstash.io",
  token:
    "AXcHASQgNTQwY2U3YWQtMmM3NC00ZTE1LTgyZDEtOWQ3NTdkNDBmNjEyYzcxMWFhMTVhNmE2NDI3ZTgxMzJiN2NmYzYwNGE1ODk=",
});

const analytics = new Redis({
  url: "https://eu2-causal-shepherd-30502.upstash.io",
  token:
    "AXcmASQgODZhMWExZTEtOWYyZC00YmM1LWE3M2YtMjMwYjYwN2E2MGVjNjUxNWNmYWE1MGMxNDUxMTliYzE0NTVjMmYxNGZiYjU=",
});

router.get("/bytecrowd/:bytecrowd", async ({ req, res }) => {
  const bytecrowd = await bytecrowds.hgetall(req.params.bytecrowd);

  if (bytecrowd !== null) res.body = bytecrowd;
  else res.body = {};
});

router.post("/update", async ({ req, res }) => {
  let name, text, language;
  name = req.body.name;
  text = req.body.text;
  language = req.body.language;

  const storedBytecrowd = await bytecrowds.hgetall(name);
  if (!storedBytecrowd)
    // If the bytecrowd doesn't exist, create it.
    await bytecrowds.hmset(name, { text: text, language: "javascript" });
  else if (
    // If at least one element(text/language) changed , update the bytecrowd.
    storedBytecrowd.text !== text ||
    storedBytecrowd.language !== language
  ) {
    // If the request doesn't contain a new text/language, use the current one.
    if (!text) text = storedBytecrowd.text;
    if (!language) language = storedBytecrowd.language;
    await bytecrowds.hmset(name, { text: text, language: language });
  }
});

router.post("/analytics", async ({ req, res }) => {
  const _updateArray = (name, stat) => {
    const array = storedDayStat[name];

    let didUpdate = false;
    if (!array.includes(stat)) {
      didUpdate = true;
      array.push(stat);
    }
    return {
      updatedArray: array,
      didUpdate: didUpdate,
    };
  };

  const page = req.body.page;
  const { country, continent } = req.cf;
  const requestIP = req.headers.get("CF-Connecting-IP");

  const _date = new Date();
  const date =
    _date.getFullYear() + " " + (_date.getMonth() + 1) + " " + _date.getDate();

  const storedDayStat = await analytics.hgetall(date);
  if (!storedDayStat) {
    // If this day wasn't recorded, create a new entry for it.
    await analytics.hmset(date, {
      hits: 1,
      addresses: [requestIP],
      uniqueVisitors: 1,
      countries: [country],
      continents: [continent],
      pages: [page],
    });
  } else {
    let { updatedArray, didUpdate } = _updateArray("addresses", requestIP);
    let uniqueVisitors = storedDayStat.uniqueVisitors;
    // If the addresses vector did update, it means a new IP visited the site.
    if (didUpdate) uniqueVisitors++;

    await analytics.hmset(date, {
      hits: storedDayStat.hits + 1,
      addresses: updatedArray,
      uniqueVisitors: uniqueVisitors,
      countries: _updateArray("countries", country).updatedArray,
      continents: _updateArray("continents", continent).updatedArray,
      pages: _updateArray("pages", page).updatedArray,
    });
  }
});

export default {
  async fetch(request, env) {
    return router.handle(env, request);
  },
};
