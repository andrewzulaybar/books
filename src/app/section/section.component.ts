import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-section',
  templateUrl: './section.component.html',
  styleUrls: ['./section.component.css'],
})
export class SectionComponent implements OnInit {
  @Input() public title: string;
  @Input() public options: string;

  public id: string;

  constructor() {}

  ngOnInit(): void {
    this.id = this.title.toLowerCase();
  }
}
