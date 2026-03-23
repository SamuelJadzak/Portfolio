import { Component } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { inject } from "@angular/core";
import { BlogApiService } from "./blog-api.service";
import { AsyncPipe } from "@angular/common";
import { CommonModule } from "@angular/common";
import { Observer, map } from "rxjs";

@Component({
  selector: "app-blog",
  imports: [CommonModule],
  templateUrl: "./blog.component.html",
  styleUrl: "./blog.component.scss",
})
export class BlogComponent {
  readonly blogService = inject(BlogApiService);
  public readonly postsMetadata$ = this.blogService.postsMetaData$;

  constructor() {
    const observer: Observer<any> = {
      next: (value) => console.log(value),
      error: (error) => console.error(error),
      complete: () => console.log("Completed"),
    };
    this.postsMetadata$.subscribe(observer);
  }
}
