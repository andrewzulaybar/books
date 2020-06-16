import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { environment } from '@env/environment';
import { IPublication } from '@models/publication.model';

@Injectable({
  providedIn: 'root',
})
export class PublicationService {
  private publicationUrl = `${environment.apiUrl}/publication`;

  constructor(private http: HttpClient) {}

  public getPublications(): Observable<IPublication[]> {
    return this.http
      .get<IPublication[]>(this.publicationUrl)
      .pipe(map((publications: IPublication[]) => publications));
  }
}
