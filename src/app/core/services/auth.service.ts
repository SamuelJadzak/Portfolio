import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable, Observer, map } from "rxjs";

interface LoginResponse {
  access_token: string;
}
type RefreshResponse = {
  access_token: string;
};

@Injectable({
  providedIn: "root",
})
export class AuthService {
  accessToken: string = "";

  constructor(private http: HttpClient) {}

  login(username: string, password: string) {
    this.http
      .post<LoginResponse>("/api/login", { username, password })
      .subscribe((res) => {
        this.accessToken = res.access_token;
        if (res.access_token) {
          console.log("Login successful", res);
        }
      });
  }

  refresh() {
    return this.http.get<RefreshResponse>("/api/refresh").pipe(
      map((res) => {
        return res.access_token;
      }),
    );
  }
}
