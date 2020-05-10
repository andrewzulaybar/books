import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  public books: { image: string; title: string; subtitle: string }[] = [
    { image: '', title: 'Title', subtitle: 'Author' },
    { image: '', title: 'Title', subtitle: 'Author' },
    { image: '', title: 'Title', subtitle: 'Author' },
    { image: '', title: 'Title', subtitle: 'Author' },
    { image: '', title: 'Title', subtitle: 'Author' },
    { image: '', title: 'Title', subtitle: 'Author' },
  ];

  public trendingOptions: string[] = ['International', 'Canada', 'My Network'];

  constructor() {}

  ngOnInit(): void {}
}
