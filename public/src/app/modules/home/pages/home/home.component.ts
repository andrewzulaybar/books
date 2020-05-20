import { Component, OnInit } from '@angular/core';
import { environment } from '@env/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  public publications: object[];
  public trendingOptions: string[] = ['International', 'Canada', 'My Network'];

  constructor() {}

  async ngOnInit() {
    try {
      this.publications = await this.fetchPublications();
    } catch (error) {
      this.publications = [];
    }
  }

  async fetchPublications(): Promise<object[]> {
    try {
      const response = await fetch(`${environment.apiUrl}/publications`);
      return response.json();
    } catch (error) {
      throw error;
    }
  }
}
