import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-card-small',
  templateUrl: './small.component.html',
  styleUrls: ['./small.component.css'],
})
export class CardSmallComponent implements OnInit {
  @Input() public content: { image: string; title: string; subtitle: string };

  constructor() {}

  ngOnInit(): void {}
}
