import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-card-small',
  templateUrl: './small.component.html',
  styleUrls: ['./small.component.css'],
})
export class CardSmallComponent implements OnInit {
  @Input() public content: { imageUrl: string; title: string; author: string };

  constructor() {}

  ngOnInit(): void {}
}
