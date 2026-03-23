import { Component } from "@angular/core";
import { RouterModule } from "@angular/router";
import { RouterOutlet } from "@angular/router";
import { RouterLink } from "@angular/router";

@Component({
  selector: "app-home",
  imports: [RouterModule, RouterOutlet, RouterLink],
  templateUrl: "./home.component.html",
  styleUrl: "./home.component.scss",
})
export class HomeComponent {}
