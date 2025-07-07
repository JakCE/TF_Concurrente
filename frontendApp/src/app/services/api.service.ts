import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private readonly baseUrl = environment.apiUrl; // o el host de tu API
  // private readonly baseUrl = 'http://localhost:8080'; // o el host de tu API

  constructor(private http: HttpClient) {}

  predecir(data: any): Observable<any> {
    return this.http.post(`${this.baseUrl}/predecir`, data);
  }

  obtenerDiccionarios(): Observable<any> {
    return this.http.get(`${this.baseUrl}/diccionarios`);
  }

  obtenerEstadisticas(): Observable<any> {
    return this.http.get(`${this.baseUrl}/estadisticas`);
  }
}