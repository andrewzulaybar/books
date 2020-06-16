import { Component, OnInit } from '@angular/core';

import { IPublication } from '@models/publication.model';
import { PublicationService } from '@service/publication.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
  providers: [PublicationService],
})
export class HomeComponent implements OnInit {
  public publications: object[];
  public trendingOptions: string[] = ['International', 'Canada', 'My Network'];

  constructor(private publicationService: PublicationService) {}

  async ngOnInit() {
    this.getPublications();
  }

  private getPublications() {
    this.publicationService
      .getPublications()
      .subscribe((publications: IPublication[]) => {
        this.publications = publications.map((publication: IPublication) => {
          const {
            imageUrl,
            work: {
              author: { firstName, lastName },
              title,
            },
          } = publication;
          const authorName = `${firstName} ${lastName}`;
          return { authorName, imageUrl, title };
        });
      });
  }
}
