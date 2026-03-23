import { Injectable, inject } from "@angular/core";
import { Observable, switchMap, shareReplay, of, tap } from "rxjs";
import { HttpClient } from "@angular/common/http";

interface PostMetaData {
  id: number;
  title: string;
}

interface Post extends PostMetaData {
  body: string[];
}

@Injectable({ providedIn: "root" })
export class BlogApiService {
  private postDataCache = new Map<number, Post>();
  private http = inject(HttpClient);
  public postsMetaData$ = this.getPostsMetaData().pipe(
    shareReplay({ bufferSize: 1, refCount: true }),
  );

  getSinglePost(id: number): Observable<Post> {
    const cachedData = this.postDataCache.get(id);
    if (cachedData) {
      return of(cachedData);
    }
    return this.http
      .get<Post>(`/api/posts/${id}`)
      .pipe(tap((response: any) => this.postDataCache.set(id, response)));
  }

  getPostsMetaData(): Observable<PostMetaData> {
    return this.http.get<PostMetaData>(`/api/posts`);
  }
}
