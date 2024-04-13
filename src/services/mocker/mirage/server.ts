import {
  RestSerializer,
  createServer,
  Factory,
  // trait
} from "miragejs";

import { createRoutes } from "./routes";
import * as models from "./models";
import { TMockServer } from "./types";

export type TRunMirageServerConfig = {
  environment?: string;
  logging?: boolean;
};

export function runServer(config: TRunMirageServerConfig = {}): TMockServer {
  return createServer({
    logging: config.logging || true,
    environment: config?.environment || "development",
    serializers: {
      reminder: RestSerializer.extend({
        include: ["list"],
        embed: true,
      }),
    },

    // models: {
    //   list: Model.extend({
    //     reminders: hasMany(),
    //   }),

    //   reminder: Model.extend({
    //     list: belongsTo(),
    //   }),
    // },

    models,

    factories: {
      list: Factory.extend({
        name(i) {
          return `List ${i}`;
        },

        // withReminders: trait({
        //   afterCreate(list, server) {
        //     if (!list.reminders.length) {
        //       server.createList("reminder", 5, { list });
        //     }
        //   }
        // })
      }),

      reminder: Factory.extend({
        text(i) {
          return `Reminder ${i}`;
        },
      }),
    },

    seeds(server) {
      server.create("reminder", { text: "Walk the dog" });
      server.create("reminder", { text: "Take out the trash" });
      server.create("reminder", { text: "Work out" });

      server.create("list", {
        name: "Home",
        reminders: [server.create("reminder", { text: "Do taxes" })],
      });

      server.create("list", {
        name: "Work",
        reminders: [server.create("reminder", { text: "Visit bank" })],
      });
    },

    routes() {
      createRoutes.call(this);

      this.get("/api/lists", (schema, request) => {
        return schema.all("list");
      });

      // this.get("/api/lists/:id/reminders", (schema, request) => {
      //   let reminders =
      //     schema.findBy("list", {
      //       id: request.params.id,
      //     })?.reminders || [];

      //   return reminders;
      // });

      this.get("/api/reminders", (schema) => {
        return schema.all("reminder");
      });

      this.post("/api/reminders", (schema, request) => {
        let attrs = JSON.parse(request.requestBody);

        return schema.create("reminder", attrs);
      });

      this.post("/api/lists", (schema, request) => {
        let attrs = JSON.parse(request.requestBody);

        return schema.create("list", attrs);
      });

      this.delete("/api/reminders/:id", (schema, request) => {
        let id = request.params.id;

        return schema.findBy("reminder", { id })?.destroy() || null;
      });

      // this.delete("/api/lists/:id", (schema, request) => {
      //   let id = request.params.id;
      //   let list = schema.findBy("list", { id });

      //   if (list?.reminders) {
      //     list?.reminders?.destroy();
      //   }

      //   return list?.destroy() || null;
      // });

      this.passthrough();
    },
  });
}
