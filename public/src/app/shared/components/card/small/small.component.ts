import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-card-small',
  templateUrl: './small.component.html',
  styleUrls: ['./small.component.css'],
})
export class CardSmallComponent implements OnInit {
  @Input() public imageUrl: string;
  @Input() public subtitle: string;
  @Input() public title: string;

  constructor() {}

  ngOnInit(): void {}
}
