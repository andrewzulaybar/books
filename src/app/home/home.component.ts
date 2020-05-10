import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  public books: { image: string; title: string; subtitle: string }[] = [
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/81X4R7QhFkL.jpg',
      title: 'Normal People',
      subtitle: 'Sally Rooney',
    },
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/91twTG-CQ8L.jpg',
      title: 'Little Fires Everywhere',
      subtitle: 'Celeste Ng',
    },
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/51j5p18mJNL.jpg',
      title: 'Where the Crawdads Sing',
      subtitle: 'Delia Owens',
    },
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/81af+MCATTL.jpg',
      title: 'The Great Gatsby',
      subtitle: 'F. Scott Fitzgerald',
    },
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/81iVsj91eQL.jpg',
      title: 'American Dirt',
      subtitle: 'Jeanine Cummins',
    },
    {
      image: 'https://images-na.ssl-images-amazon.com/images/I/91Xq+S+F2jL.jpg',
      title: 'Atomic Habits',
      subtitle: 'James Clear',
    },
  ];

  public trendingOptions: string[] = ['International', 'Canada', 'My Network'];

  constructor() {}

  ngOnInit(): void {}
}
