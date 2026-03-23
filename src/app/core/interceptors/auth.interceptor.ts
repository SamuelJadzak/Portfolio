import { inject } from "@angular/core";
import { HttpInterceptorFn, HttpRequest } from "@angular/common/http";
import { AuthService } from "../services/auth.service";
import { environment } from "../../../environments/environment";
import { Observable, switchMap, tap, timeout } from "rxjs";
import { HttpEvent } from "@angular/common/http";

export const authInterceptor: HttpInterceptorFn = (
  request,
  next,
): Observable<HttpEvent<unknown>> => {
  const authService = inject(AuthService);
  switch (request.url) {
    case "/api/login":
    case "/api/refresh":
      const requestWithBaseUrl = setBaseUrl(request, environment.baseServerUrl);
      const requestWithCredentials = withCredentials(requestWithBaseUrl);
      return next(requestWithCredentials);
    default:
      return authService.refresh().pipe(
        switchMap((newToken) => {
          const requestWithBaseUrl = setBaseUrl(
            request,
            environment.baseServerUrl,
          );
          const requestWithAuth = setAuthorizationHeader(
            requestWithBaseUrl,
            newToken,
          );
          return next(requestWithAuth);
        }),
        timeout(2000),
        tap(() => {
          console.log("Request intercepted and modified with auth token");
        }),
      );
  }
};

function setBaseUrl(
  request: HttpRequest<any>,
  baseServerUrl: string,
): HttpRequest<any> {
  return request.clone({
    url: `${baseServerUrl}${request.url}`,
  });
}

function setAuthorizationHeader(
  request: HttpRequest<any>,
  authToken: string,
): HttpRequest<any> {
  return request.clone({
    setHeaders: { Authorization: `Bearer ${authToken}` },
  });
}

function withCredentials(request: HttpRequest<any>): HttpRequest<any> {
  return request.clone({
    withCredentials: true,
  });
}
