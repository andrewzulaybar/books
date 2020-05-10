import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-card-small',
  templateUrl: './card-small.component.html',
  styleUrls: ['./card-small.component.css'],
})
export class CardSmallComponent implements OnInit {
  @Input() public content: { image: string; title: string; subtitle: string };

  constructor() {}

  ngOnInit(): void {}
}
