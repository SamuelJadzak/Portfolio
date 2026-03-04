import { Component } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { inject } from "@angular/core";
import { BlogApiService } from "./blog-api.service";
import { AsyncPipe } from "@angular/common";
@Component({
  selector: "app-blog",
  imports: [],
  templateUrl: "./blog.component.html",
  styleUrl: "./blog.component.scss",
})
export class BlogComponent {
  readonly blogService = inject(BlogApiService);
  public readonly postsMetadata$ = this.blogService.postsMetaData$;
}
