import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-section-container',
  templateUrl: './container.component.html',
  styleUrls: ['./container.component.css'],
})
export class SectionContainerComponent implements OnInit {
  @Input() public infoRoute: string;
  @Input() public options: string;
  @Input() public title: string;

  constructor() {}

  ngOnInit(): void {}
}
