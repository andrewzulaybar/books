import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-section-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css'],
})
export class SectionHeaderComponent implements OnInit {
  @Input() public options: string[];
  @Input() public title: string;

  constructor() {}

  ngOnInit(): void {}
}
