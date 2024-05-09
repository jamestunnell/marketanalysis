import UniversalRouter from "universal-router";
import context from "../context";

const routes = [
  {
    path: "",
    action: async (context) => {
      const { default: page } = await import("./home");
      return page(context);
    },
  },
  {
    path: "securities",
    action: async (context) => {
      const { default: page } = await import("./securities");
      return page(context);
    },
  },
  {
    path: "graphs",
    action: async (context) => {
      const { default: page } = await import("./graphs");
      return page(context);
    },
  },

  {
    path: "(.*)",
    action: () => "Nothing there",
  },
];

const router = new UniversalRouter(routes, { context });
export default router;