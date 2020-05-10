import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  public books: number[] = Array(6).fill(0);

  public trendingOptions: string[] = ['International', 'Canada', 'My Network'];

  constructor() {}

  ngOnInit(): void {}
}
