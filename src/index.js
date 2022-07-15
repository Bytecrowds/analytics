import Router from "@tsndr/cloudflare-worker-router";
import { Redis } from "@upstash/redis/cloudflare";

const router = new Router();

// Because cloudflare workers expose the env on fetch, we need to manually set the vars
const redis = new Redis({
  url: "https://eu2-large-spider-30471.upstash.io",
  token:
    "AXcHASQgZjU1M2Q2MjYtNzdhNy00NjY1LWIwZTktMDY1ZjM0MzFmNGQ0ZjNhMGQ4ZWI0MDA3NGY1OTg4YjU4NjQzN2E2ZDA3NzY=",
});

router.get("/bytecrowd/:bytecrowd", async ({ req, res }) => {
  const bytecrowd = await redis.hgetall(req.params.bytecrowd);

  if (bytecrowd !== null) res.body = bytecrowd;
  else res.body = {};
});

router.post("/update", async ({ req, res }) => {
  let name, text, language;
  name = req.body.name;
  text = req.body.text;
  language = req.body.language;

  const storedBytecrowd = await redis.hgetall(name);
  if (!storedBytecrowd)
    // if the bytecrowd doesn't exist, create it
    await redis.hmset(name, { text: text, language: "javascript" });
  else if (
    // if at least one element(text/language) changed , update the bytecrowd
    storedBytecrowd.text !== text ||
    storedBytecrowd.language !== language
  ) {
    // if the request doesn't contain a new text/language, use the current one
    if (!text) text = storedBytecrowd.text;
    if (!language) language = storedBytecrowd.language;
    await redis.hmset(name, { text: text, language: language });
  }
});

export default {
  async fetch(request, env) {
    return router.handle(env, request);
  },
};
