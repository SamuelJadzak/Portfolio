import { Routes } from "@angular/router";
import { HomeComponent } from "./home/home.component";
import { BlogComponent } from "./blog/blog.component";
export const routes: Routes = [
  {
    path: "",
    children: [
      {
        path: "",
        component: HomeComponent,
        data: {
          title: "Home",
          icon: "home",
        },
      },
      {
        path: "blog/",
        component: BlogComponent,
        data: {
          title: "Blog",
          icon: "blog",
        },
      },
    ],
  },
];
