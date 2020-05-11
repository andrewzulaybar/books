import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NavigationHeaderComponent } from './header.component';

describe('NavigationBarComponent', () => {
  let component: NavigationHeaderComponent;
  let fixture: ComponentFixture<NavigationHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [NavigationHeaderComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NavigationHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
