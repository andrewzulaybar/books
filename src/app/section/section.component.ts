import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-section',
  templateUrl: './section.component.html',
  styleUrls: ['./section.component.css'],
})
export class SectionComponent implements OnInit {
  @Input() public infoRoute: string;
  @Input() public options: string;
  @Input() public title: string;

  constructor() {}

  ngOnInit(): void {}
}
